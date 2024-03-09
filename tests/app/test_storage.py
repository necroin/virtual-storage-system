from conftest import Environment
import json

def test_fylesystem(environment: Environment):
    token = environment.project.server.token()
    response = environment.project.server.post_request("/{}/storage/filesystem".format(token), data="/")
    filesystem = json.loads(response)
    print(filesystem)
    assert False
    