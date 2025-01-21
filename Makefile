add_task:
	curl -X POST -H "Content-Type: application/json" \
		-d '{"type": "send_email", "payload": "example@example.com"}' \
		http://localhost:8080/tasks

