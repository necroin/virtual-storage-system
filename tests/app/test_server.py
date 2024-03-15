from conftest import Environment


def test_status(environment: Environment):
    assert environment.project.server.status() == "OK"
