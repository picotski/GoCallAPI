FROM python:3.12-bullseye

WORKDIR /test

COPY tests/requirements.txt ./

RUN pip install -r requirements.txt

COPY tests ./

# ENTRYPOINT [ "sleep", "10000" ]
ENTRYPOINT [ "pytest" ]
CMD [ "./call_test.py" ]