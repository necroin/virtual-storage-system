import requests

class Server:
    def __init__(self, url) -> None:
        self.url = url

    def status(self):
        response = requests.get(self.url+"/status", verify=False)
        response.raise_for_status()
        return response.text

    def post_request(self, path, data):
        response = requests.post(self.url + path, json=data, verify=False)
        return response.text
    
    def get_request(self, path, data):
        response = requests.get(self.url + path, json=data, verify=False)
        return response.text