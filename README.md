Clock Application


Project Description

This application is a clock that prints "tick" every second, "tock" every minute, and "bong" every hour. The messages can be updated in real-time via a REST API. The application is designed to run for three hours before automatically terminating.


User Story

As a user, I want a clock application that prints messages at regular intervals: "tick" every second, "tock" every minute, and "bong" every hour. I should be able to change these messages at any time while the application is running, without needing to stop it, to customize the output according to my needs. Additionally, I want the application to log each signal the clock triggers into a database to track events historically.


Design Decisions

I used Go goroutines to handle concurrency between the clock and the HTTP server.
I chose to use a REST API to allow message updates, providing flexibility and ease of use.
Each signal triggered by the clock is stored in an SQLite database.


Technologies Used

Go: For the main logic of the application.
net/http: To implement the HTTP server.
SQLite: To store records of message changes.


Challenges and Solutions

Concurrency Handling: Used channels to communicate message changes between the HTTP server and the clock.
Data Persistence: Opted for SQLite due to its simplicity and efficiency for a project of this size.


Execution Instructions

Clone the repository
Navigate to the project directory: cd go-clock-app
Run the program: go run main.go


How to Test the Application

Update Messages

To update the messages, use the following cURL command:
curl -X POST -d '{"TickMessage":"quack", "TockMessage":"took", "BongMessage":"bang"}' -H "Content-Type: application/json" http://localhost:8080/update-signals

Update Intervals

An additional endpoint has been created to set the initial interval values for the messages. Use the following cURL command to update the intervals:
curl -X POST -d '{"TickInterval":"3s", "TockInterval":"10s", "BongInterval":"1m"}' -H "Content-Type: application/json" http://localhost:8080/update-intervals
