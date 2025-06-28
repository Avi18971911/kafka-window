#!/bin/bash

# Kafka Setup Script - Creates topics, messages, and consumer offsets

echo "🚀 Setting up Kafka topics and test data..."

# Wait for Kafka to be ready
echo "⏳ Waiting for Kafka to be ready..."
docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1
while [ $? -ne 0 ]; do
    sleep 2
    docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1
done
echo "✅ Kafka is ready!"

# Create test topics
echo "📝 Creating topics..."
docker exec kafka kafka-topics --create --topic user-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
docker exec kafka kafka-topics --create --topic order-events --bootstrap-server localhost:9092 --partitions 2 --replication-factor 1
docker exec kafka kafka-topics --create --topic logs --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
docker exec kafka kafka-topics --create --topic notifications --bootstrap-server localhost:9092 --partitions 4 --replication-factor 1

echo "📊 Topics created:"
docker exec kafka kafka-topics --list --bootstrap-server localhost:9092

# Function to produce messages to a topic
produce_messages() {
    local topic=$1
    local messages=$2
    echo "📤 Producing messages to $topic..."
    echo "$messages" | docker exec -i kafka kafka-console-producer --topic $topic --bootstrap-server localhost:9092
}

# Produce sample messages
echo "🎯 Adding sample messages..."

# User events
produce_messages "user-events" "$(cat << 'EOF'
{"user_id": "user-001", "event": "login", "timestamp": "2025-06-28T10:00:00Z"}
{"user_id": "user-002", "event": "signup", "timestamp": "2025-06-28T10:05:00Z"}
{"user_id": "user-001", "event": "page_view", "page": "/dashboard", "timestamp": "2025-06-28T10:10:00Z"}
{"user_id": "user-003", "event": "login", "timestamp": "2025-06-28T10:15:00Z"}
{"user_id": "user-002", "event": "purchase", "amount": 49.99, "timestamp": "2025-06-28T10:20:00Z"}
EOF
)"

# Order events
produce_messages "order-events" "$(cat << 'EOF'
{"order_id": "ord-001", "user_id": "user-002", "status": "created", "total": 49.99}
{"order_id": "ord-001", "user_id": "user-002", "status": "paid", "total": 49.99}
{"order_id": "ord-002", "user_id": "user-001", "status": "created", "total": 129.50}
{"order_id": "ord-001", "user_id": "user-002", "status": "shipped", "tracking": "TRK123"}
EOF
)"

# Log messages
produce_messages "logs" "$(cat << 'EOF'
[INFO] Application started successfully
[WARN] High memory usage detected: 85%
[ERROR] Database connection timeout after 30s
[INFO] User user-001 logged in from IP 192.168.1.100
[DEBUG] Processing batch of 150 events
[ERROR] Failed to send notification to user-003
EOF
)"

# Notifications
produce_messages "notifications" "$(cat << 'EOF'
{"type": "email", "recipient": "user-001@example.com", "subject": "Welcome!", "sent": true}
{"type": "sms", "recipient": "+1234567890", "message": "Your order has shipped", "sent": true}
{"type": "push", "recipient": "user-002", "title": "New message", "sent": false, "error": "Device not registered"}
{"type": "email", "recipient": "user-003@example.com", "subject": "Password reset", "sent": true}
EOF
)"

echo "✅ Messages produced to all topics!"

# Create consumers to generate consumer offsets
echo "🔄 Creating consumer groups to populate __consumer_offsets..."

# Consumer group 1: Read from user-events
echo "👥 Starting consumer group 'analytics-team'..."
timeout 10s docker exec kafka kafka-console-consumer \
    --topic user-events \
    --bootstrap-server localhost:9092 \
    --group analytics-team \
    --from-beginning > /dev/null 2>&1 &

# Consumer group 2: Read from order-events
echo "👥 Starting consumer group 'order-processors'..."
timeout 10s docker exec kafka kafka-console-consumer \
    --topic order-events \
    --bootstrap-server localhost:9092 \
    --group order-processors \
    --from-beginning > /dev/null 2>&1 &

# Consumer group 3: Read from logs
echo "👥 Starting consumer group 'log-aggregators'..."
timeout 10s docker exec kafka kafka-console-consumer \
    --topic logs \
    --bootstrap-server localhost:9092 \
    --group log-aggregators \
    --from-beginning > /dev/null 2>&1 &

# Consumer group 4: Read from multiple topics
echo "👥 Starting consumer group 'monitoring-service'..."
timeout 10s docker exec kafka kafka-console-consumer \
    --topic notifications \
    --bootstrap-server localhost:9092 \
    --group monitoring-service \
    --from-beginning > /dev/null 2>&1 &

# Wait for consumers to commit offsets
echo "⏳ Waiting for consumer offset commits..."
sleep 15

# Show results
echo "🎉 Setup complete! Here's what we created:"
echo ""
echo "📋 Topics:"
docker exec kafka kafka-topics --list --bootstrap-server localhost:9092

echo ""
echo "👥 Consumer Groups:"
docker exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --list

echo ""
echo "📊 Consumer Group Details:"
for group in analytics-team order-processors log-aggregators monitoring-service; do
    echo "--- $group ---"
    docker exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group $group 2>/dev/null || echo "No offsets yet for $group"
    echo ""
done

echo "🔍 __consumer_offsets topic info:"
docker exec kafka kafka-topics --describe --topic __consumer_offsets --bootstrap-server localhost:9092

echo ""
echo "✨ All done! Your Kafka cluster now has:"
echo "   • 4 custom topics with sample data"
echo "   • 4 consumer groups with committed offsets"
echo "   • __consumer_offsets topic populated with offset data"
echo ""
echo "💡 You can now explore the data using Kafka UI at http://localhost:8080"