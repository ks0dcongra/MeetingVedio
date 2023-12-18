# 自定義啟動腳本的內容
echo "Running custom entrypoint script..."

# 啟動 MongoDB
exec mongod --auth --bind_ip_all

chmod +x docker-entrypoint.sh

print("Executing mongo-init.js script...");

mongo -- "$MONGO_INITDB_DATABASE" <<EOF
db = db.getSiblingDB('admin')
db.auth('$MONGO_INITDB_ROOT_USERNAME', '$MONGO_INITDB_ROOT_PASSWORD')
db = db.getSiblingDB('$MONGO_INITDB_DATABASE')
db.createUser({
  user: "$MONGO_USERNAME",
  pwd: "$MONGO_PASSWORD",
  roles: [
  { role: 'readWrite', db: '$MONGO_INITDB_DATABASE' }
  ]
})
EOF


# print("Executing mongo-init.js script...");

# mongo -- "$MONGO_INITDB_DATABASE" -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" <<EOF
# db.createUser({
#   user: "$MONGO_USERNAME",
#   pwd: "$MONGO_PASSWORD",
#   roles: [
#     { role: 'readWrite', db: '$MONGO_INITDB_DATABASE' }
#   ]
# })
# EOF