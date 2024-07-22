import os
import requests
import pytest
from db.db_services import clear_db

base_url = os.getenv('TEST_HOST_ADDR') or 'http://localhost:8010'

@pytest.fixture
def setup_before_get_call_test():
  clear_db()

  return 1

# Verify that the server is on
def test_if_server_health_check_isvalid():
  # Arrange
  url = f'{base_url}/health'

  # Act
  res = requests.get(url)

  # Assert
  assert res.status_code == 200

# Get a call that exists
def test_get_call_valid(setup_before_get_call_test):
  # Arrange
  url = f'{base_url}/call'
  obj = {
    "caller": "john",
    "recipient": "pierre"
  }

  # Act
  ## Create a call
  post_res = requests.post(url, json=obj)

  id = str(post_res.json()['id'])

  ## Get the call
  res = requests.get(f'{url}/{id}')

  # Assert
  assert res.status_code == 200

def test_get_call_invalid():
  # Arrange
  url = f'{base_url}/call/1'

  # Act
  res = requests.get(url)

  # Assert
  assert res.status_code == 404

def test_create_call_valid():
  # Arrange
  url = f'{base_url}/call'
  obj = {
    "caller": "john",
    "recipient": "pierre"
  }

  # Act
  res = requests.post(url, json=obj)

  # Assert
  assert res.status_code == 201

def test_create_call_invalid():
  # Arrange
  url = f'{base_url}/call'

  # Act
  res = requests.post(url)

  # Assert
  assert res.status_code == 400

def test_delete_call_valid():
  # Arrange
  url = f'{base_url}/call/1'

  # Act
  res = requests.delete(url)

  # Assert
  assert res.status_code == 200

def test_stop_call_valid():
  # Arrange
  create_url = f'{base_url}/call'
  stop_url = f'{base_url}/stop'
  obj = {
    "caller": "john",
    "recipient": "pierre"
  }

  # Act
  res = requests.post(create_url, json=obj)

  id = str(res.json()['id'])

  stop_res = requests.get(f'{stop_url}/{id}')

  # Assert
  assert stop_res.status_code == 200
  assert stop_res.json()['status'] == 'Ended'

def test_stop_call_invalid_not_found():
  # Arrange
  stop_url = f'{base_url}/stop'

  # Act
  stop_res = requests.get(f'{stop_url}/1')

  # Assert
  assert stop_res.status_code == 404

def test_stop_call_invalid_already_ended():
  # Arrange
  create_url = f'{base_url}/call'
  stop_url = f'{base_url}/stop'
  obj = {
    "caller": "john",
    "recipient": "pierre"
  }

  # Act
  res = requests.post(create_url, json=obj)

  id = str(res.json()['id'])

  requests.get(f'{stop_url}/{id}')
  stop_res = requests.get(f'{stop_url}/{id}')

  # Assert
  assert stop_res.status_code == 400