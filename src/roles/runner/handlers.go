package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"vss/src/logger"
	"vss/src/message"
	"vss/src/settings"
	"vss/src/utils"

	"github.com/gorilla/mux"
	"github.com/lxn/win"
	"github.com/necroin/golibs/utils/winapi"
	"github.com/necroin/golibs/utils/winutils"
	"github.com/necroin/golibs/winappstream"
)

func (runner *Runner) OpenFileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	openResponse := &message.OpenResponse{}
	defer json.NewEncoder(responseWriter).Encode(openResponse)

	openRequest := &message.OpenRequest{}
	if err := json.NewDecoder(request.Body).Decode(openRequest); err != nil {
		openResponse.Error = err
		logger.Error("[Runner] [OpenFileHandler] failed decode message: %s", err)
		return
	}

	openPath := openRequest.Path
	selfUrl := fmt.Sprintf("%s/%s", runner.config.Url, runner.config.User.Token)
	if selfUrl != openRequest.SrcUrl {
		openPath = fmt.Sprintf("./tmp/%s", filepath.Base(openRequest.Path))

		copyRequest := &message.CopyRequest{
			OldPath: openRequest.Path,
			NewPath: openPath,
			SrcUrl:  openRequest.SrcUrl,
		}
		runner.connector.SendPostRequest(fmt.Sprintf("%s/storage/copy/%s", selfUrl, openRequest.Type), copyRequest)
	}

	logger.Info("[Runner] [OpenFileHandler] open %s", openPath)

	execTool, execArgs := runner.GetRunCommand(path.Clean(openPath))

	cmd := exec.Command(execTool, execArgs...)
	if err := cmd.Start(); err != nil {
		openResponse.Error = err
		logger.Error("[Runner] [OpenFileHandler] failed start process: %s", err)
		return
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			logger.Error("[Runner] [OpenFileHandler] failed finish process by wait: %s", err)
		}
	}()
	logger.Debug("[Runner] [OpenFileHandler] process started with pid: %d", cmd.Process.Pid)
	openResponse.Pid = cmd.Process.Pid
}

func (runner *Runner) AppStreamHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	pid, _ := strconv.Atoi(params["pid"])
	appStreamPid := winapi.ProcessId(0)
	procName := ""

	var childProcesses []*winutils.Process
	if err := utils.Try(func() error {
		allProcesses, err := winutils.GetAllProcesses()
		if err != nil {
			return fmt.Errorf("[Runner] [AppStreamHandler] failed get all processes")
		}

		childProcesses = winutils.FindProcessesByParentPid(allProcesses, winapi.ProcessId(pid))
		if len(childProcesses) == 0 {
			return fmt.Errorf("[Runner] [AppStreamHandler] failed find child processes for pid %d", pid)
		}
		logger.Debug("[Runner] [AppStreamHandler] childs of %d pid: %s", pid, childProcesses)

		for _, childProcess := range childProcesses {
			procName = childProcess.Executable
			appPid := childProcess.Pid
			logger.Debug("[Runner] [AppStreamHandler] app process started with pid: %d", appPid)

			rect, ok := utils.FindValidRect(appPid)
			logger.Debug("[Runner] [AppStreamHandler] valid rect: %v", rect)
			if ok {
				appStreamPid = appPid
				return nil
			}
		}

		return fmt.Errorf("[Runner] [AppStreamHandler] failed find valid rect")
	}, 10, 1*time.Second); err != nil {
		logger.Error(err.Error())
	}

	if appStreamPid == 0 {
		allProcesses, err := winutils.GetAllProcesses()
		if err != nil {
			err = fmt.Errorf("[Runner] [AppStreamHandler] failed get all processes and find app: %s", err)
			logger.Error(err.Error())
			responseWriter.Write([]byte(err.Error()))
			return
		}

		for _, process := range allProcesses {
			if process.Executable == procName {
				_, ok := utils.FindValidRect(process.Pid)
				if ok {
					appStreamPid = process.Pid
					break
				}
			}
		}

	}

	if appStreamPid == 0 {
		err := errors.New("[Runner] [AppStreamHandler] failed find app")
		logger.Error(err.Error())
		responseWriter.Write([]byte(err.Error()))
		return
	}

	response, err := runner.connector.SendRequest(fmt.Sprintf("%s/%s/runner/stream/direct/%d", runner.config.Url, runner.config.User.Token, appStreamPid), []byte{}, "GET")
	if err != nil {
		logger.Error(err.Error())
		responseWriter.Write([]byte(err.Error()))
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err.Error())
		responseWriter.Write([]byte(err.Error()))
		return
	}
	responseWriter.Write(data)
}

func (runner *Runner) AppDirectStreamHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	pid, _ := strconv.Atoi(params["pid"])
	logger.Info("[Runner] [AppDirectStreamHandler] requested stream with pid: %d", pid)

	_, ok := runner.streamSessions[pid]
	if !ok {
		runner.sessuinMutex.Lock()
		defer runner.sessuinMutex.Unlock()

		app, err := winappstream.NewApp(winapi.ProcessId(pid))
		if err != nil {
			responseWriter.Write([]byte(fmt.Sprintf("[Runner] [AppDirectStreamHandler] failed create app stream: %s", err)))
			return
		}

		streamSession := &StreamSession{
			app:            app,
			handler:        app.HttpImageCaptureHandler(),
			lastHandleTime: time.Now(),
		}
		runner.streamSessions[pid] = streamSession

		app.LaunchStream()

		go func(pid int) {
			for {
				now := time.Now()
				if now.Sub(streamSession.lastHandleTime).Seconds() > time.Duration(30*time.Second).Seconds() {
					runner.sessuinMutex.Lock()
					defer runner.sessuinMutex.Unlock()
					streamSession.app.Destroy()
					delete(runner.streamSessions, pid)
					logger.Info("[Runner] [AppDirectStreamHandler] stream pid %d closed by expiration", pid)
					return
				}
				time.Sleep(10 * time.Second)
			}
		}(pid)

		logger.Info("[Runner] [AppDirectStreamHandler] stream for pid %d created", pid)
	}

	pageInfo := message.PageInfo{
		Url:    runner.config.Url,
		Token:  runner.config.User.Token,
		Style:  settings.GetAppStreamPageStyle(),
		Script: settings.GetAppStreamPageScript(),
		Pid:    pid,
	}
	pageTemplate, _ := template.New("StreamPage").Parse(settings.GetAppStreamPage())
	pageTemplate.Execute(responseWriter, pageInfo)
}

func (runner *Runner) AppImageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	pid, _ := strconv.Atoi(params["pid"])

	streamSession, ok := runner.streamSessions[pid]
	if !ok {
		responseWriter.Write([]byte(fmt.Sprintf("Stream of pid %d not exists", pid)))
		return
	}
	streamSession.lastHandleTime = time.Now()
	streamSession.handler.ServeHTTP(responseWriter, request)
}

func (runner *Runner) AppMouseEventHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	pid, _ := strconv.Atoi(params["pid"])

	event := &message.MouseEvent{}
	json.NewDecoder(request.Body).Decode(event)

	coords := event.Coords
	cursorPotion := win.POINT{}
	win.GetCursorPos(&cursorPotion)

	handles := winutils.GetWindowHandlesByProcessId(winapi.ProcessId(pid))

	captureRect, _ := winutils.GetCaptureRectByHandles(handles)
	coords.X += captureRect.Left
	coords.Y += captureRect.Top
	for _, handle := range handles {
		windowRect, _ := winapi.GetWindowRect(handle)
		clientRect, _ := winapi.GetClientRect(handle)
		if winutils.RectEqual(captureRect, winutils.GetCaptureRect(windowRect, clientRect)) {
			win.SetForegroundWindow(win.HWND(handle))
			win.SetActiveWindow(win.HWND(handle))
			if event.Type == "leftDown" {
				// fmt.Println("leftDown", coords)
				utils.MouseMove(coords.X, coords.Y)
				utils.MouseLeftDown()
			}
			if event.Type == "leftUp" {
				// fmt.Println("leftUp", coords)
				utils.MouseLeftUp()
			}
			if event.Type == "move" {
				// fmt.Println("move", coords)
				utils.MouseMove(coords.X, coords.Y)
			}
			if event.Type == "scroll" {
				// fmt.Println("scroll", event.Scroll)
				utils.MouseWheel(coords.X, coords.Y, event.Scroll.Y)
			}
			break
		}
	}

	// utils.MouseMove(cursorPotion.X, cursorPotion.Y)
}
