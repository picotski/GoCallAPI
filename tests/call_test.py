import os
import requests
import pytest
from db.db_services import clear_db, create_call

base_url = os.getenv('TEST_HOST_ADDR')

@pytest.fixture
def setup_before_all_test():
  clear_db()

  return 1

# Verify that the server is on
def test_server_health_check_isvalid(setup_before_all_test):
  # Arrange
  url = f'{base_url}/health'

  # Act
  res = requests.get(url)

  # Assert
  assert res.status_code == 200

# Get a call that exists
def test_get_single_call_isvalid(setup_before_all_test):
  # Arrange
  url = f'{base_url}/call'

  ## Create a call
  id = create_call()

  # Act
  ## Get the call
  res = requests.get(f'{url}/{id}')

  # Assert
  assert res.status_code == 200

def test_get_single_call_isinvalid(setup_before_all_test):
  # Arrange
  url = f'{base_url}/call/1'

  # Act
  res = requests.get(url)

  # Assert
  assert res.status_code == 404

def test_create_single_call_isvalid(setup_before_all_test):
  # Arrange
  url = f'{base_url}/call'
  obj = {
    "caller": "John",
    "recipient": "Pierre"
  }

  # Act
  res = requests.post(url, json=obj)

  # Assert
  assert res.status_code == 201
  assert res.json()['caller'] == 'John'
  assert res.json()['recipient'] == 'Pierre'
  assert res.json()['status'] == 'Ongoing'

def test_create_single_call_with_no_body_isinvalid(setup_before_all_test):
  # Arrange
  url = f'{base_url}/call'

  # Act
  res = requests.post(url)

  # Assert
  assert res.status_code == 400

def test_delete_single_call_isvalid(setup_before_all_test):
  # Arrange
  url = f'{base_url}/call/1'

  # Act
  res = requests.delete(url)

  # Assert
  assert res.status_code == 200

def test_stop_call_valid(setup_before_all_test):
  # Arrange
  stop_url = f'{base_url}/stop'

  ## Create a call
  id = create_call()

  # Act
  ## Stop call
  stop_res = requests.get(f'{stop_url}/{id}')

  # Assert
  assert stop_res.status_code == 200
  assert stop_res.json()['status'] == 'Ended'

def test_stop_call_invalid_not_found(setup_before_all_test):
  # Arrange
  stop_url = f'{base_url}/stop'

  # Act
  stop_res = requests.get(f'{stop_url}/1')

  # Assert
  assert stop_res.status_code == 404

def test_stop_call_invalid_already_ended(setup_before_all_test):
  # Arrange
  stop_url = f'{base_url}/stop'
  
  ## Create a call
  id = create_call()
  ## Stop call 1
  requests.get(f'{stop_url}/{id}')

  # Act
  ## Attempt to stop call
  stop_res = requests.get(f'{stop_url}/{id}')

  # Assert
  assert stop_res.status_code == 400