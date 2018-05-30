# Using the latest golang version 
FROM python:3.4-alpine

# Set the maintainer name 
MAINTAINER Andrew Tokarskiy <tokarskiy.a@gmail.com>

# Set the application directory 
WORKDIR /app
RUN pip install flask
RUN pip install requests

# Copy the code from current folder 
# to /app inside the container 
COPY . /app 

# Make port 80 available for links 
# and/or publish 
EXPOSE 80

# Command to run the service
CMD python app.py