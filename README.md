# Logi – Door-to-Door Delivery API

Logi is a powerful door-to-door delivery API service that enables users to send and receive packages with real-time tracking. From dispatch to drop-off, it streamlines coordination between senders, receivers, and drivers — without requiring the receiver to install an app.

---

## 🚀 Features

* 📦 Create and manage deliveries
* 🚚 Assign and track drivers in real-time
* 🧭 Monitor package status updates through each stage
* 📬 Notify senders/receivers during transit events

---

## 🛠️ Technologies Used

* **Language:** Go
* **Framework:** Gin
* **Database:** PostgreSQL, MongoDB
* **Caching:** Redis
* **Message Broker / Background Jobs:** RabbitMQ
* **Containerization:** Docker
* **Hosting:** AWS

---

## 🧪 Testing

* Supports **integration testing** and **end-to-end (e2e) testing** only
* All supported test commands are defined as **shortcuts** in the `Makefile`

---

## ⚙️ CI/CD

* CI/CD pipeline is configured to:

  * Run integration and e2e tests
  * Build and deploy Docker containers
  * Push updates to AWS automatically

---

## 📂 API Documentation

* OpenAPI documentation available at:

  ```
  /docs or /swagger
  ```
