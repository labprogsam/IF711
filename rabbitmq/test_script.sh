#!/usr/bin/zsh

RUNS=30

CLIENT_ID=$1
FILENAME="testfile_list_$CLIENT_ID.txt"

# LIST command
echo "🚀 Starting LIST test loop ($RUNS runs)..."
echo "========= LIST command test ========= " > $FILENAME
for i in $(seq 1 $RUNS); do
    echo "🧪 Test Run #$i; Cmd: \`go run ./client LIST $CLIENT_ID\`"

    # Create a unique file for this run
    echo "$(go run ./client LIST $CLIENT_ID 2>&1)" >> $FILENAME
done

# UPLOAD command
echo "🚀 Starting UPLOAD test loop ($RUNS runs)..."
echo "========= UPLOAD command test ========= " >> $FILENAME
for i in $(seq 1 $RUNS); do
    echo "🧪 Test Run #$i; Cmd: \`go run ./client UPLOAD $CLIENT_ID test.txt\`"

    # Create a unique file for this run
    echo "$(go run ./client UPLOAD $CLIENT_ID test.txt 2>&1)" >> $FILENAME
done

# DOWNLOAD command
echo "🚀 Starting DOWNLOAD test loop ($RUNS runs)..."
echo "========= DOWNLOAD command test ========= " >> $FILENAME
for i in $(seq 1 $RUNS); do
    echo "🧪 Test Run #$i; Cmd: \`go run ./client DOWNLOAD $CLIENT_ID test.txt\`"

    # Create a unique file for this run
    echo "$(go run ./client DOWNLOAD $CLIENT_ID test.txt 2>&1)" >> $FILENAME
done
