import os
import requests

base_url = os.getenv('TEST_HOST_ADDR')

# Verify that the server is on
def test_health_check_valid():
  # Arrange
  url = f'{base_url}/health'

  # Act
  res = requests.get(url)

  print(base_url)

  # Assert
  assert res.status_code == 200

# Get a call that exists
def test_get_call_valid():
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

  # Cleanup
  requests.delete(f'{url}/{id}')

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

  # Cleanup
  id = str(res.json()['id'])

  requests.delete(f'{url}/{id}')

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

  # Cleanup
  requests.delete(f'{create_url}/{id}')

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

  # Cleanup
  requests.delete(f'{create_url}/{id}')