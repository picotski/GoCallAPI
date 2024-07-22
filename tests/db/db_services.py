from configparser import ConfigParser
import psycopg2
from datetime import datetime

def load_config(filename='database.ini', section='postgresql'):
  parser = ConfigParser()
  parser.read(filename)

  # get section, default to postgresql
  config = {}
  if parser.has_section(section):
    params = parser.items(section)
    for param in params:
      config[param[0]] = param[1]
  else:
    raise Exception('Section {0} not found in the {1} file'.format(section, filename))

  return config

def connect():
    config = load_config()
    """ Connect to the PostgreSQL database server """
    try:
        # connecting to the PostgreSQL server
        with psycopg2.connect(**config) as conn:
            print('Connected to the PostgreSQL server.')
            return conn
    except (psycopg2.DatabaseError, Exception) as error:
        print(error)

def create_call():
  statement = """
    INSERT INTO calls(caller, recipient, status, start_time, end_time) 
    VALUES(%s, %s, %s, %s, %s) 
    RETURNING id
  """

  id = None

  try:
    with connect() as conn:
      with conn.cursor() as cur:
        cur.execute(statement, ('John', 'Pierre', 'Ongoing', datetime.now(), datetime.now()))

        rows = cur.fetchone()
        if rows:
          id = rows[0]
  except (psycopg2.DatabaseError, Exception) as error:
    return error

  return id

def clear_db():
  delete_statement = 'DROP TABLE calls'
  create_statement = """
    CREATE TABLE calls(
			id SERIAL PRIMARY KEY,
			caller TEXT,
			recipient TEXT,
			status TEXT,
			start_time TIMESTAMP,
			end_time TIMESTAMP
		)
  """

  try:
    with connect() as conn:
      with conn.cursor() as cur:
        cur.execute(delete_statement)
        cur.execute(create_statement)
  except (psycopg2.DatabaseError, Exception) as error:
    print(error)