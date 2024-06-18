#!/bin/bash
function handle_error() {
	local exit_code=$?
	local line_no=$1
	echo "Error occurred on line $line_no. Exit code: $exit_code"
	zenity --info --title="Error" --text="Could not save todo. Error occurred on line $line_no. Exit code: $exit_code" --icon=error
	exit $exit_code
}

# Get todo description
text=$(zenity --entry --text="" --title="TODO" --width=500 --height=100)

if [ -z "$text" ]; then
	exit
fi

text=$(echo "$text" | xargs)

createdDate=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# read current context
json=$(cat ~/.tasktango/state.json)
context=$(echo "$json" | jq -r '.activeContext')

function context_change() {
	targetContext=$1
	if [ "$context" != "$targetContext" ]; then
		zenity --question --text="Set context to $targetContext?" --default-cancel
		if [ $? -eq 0 ]; then
			context=$targetContext
			json=$(echo "$json" | jq ".activeContext = \"$context\"")
			echo "$json" >~/.tasktango/state.json
		fi
	fi
}
if [ $(date +%H) -ge 17 ] || [ $(date +%H) -le 9 ]; then
	context_change "personal"
fi

if [ $(date +%H) -ge 9 ] && [ $(date +%H) -le 16 ]; then
	context_change "work"
fi

trap 'handle_error $LINENO' ERR
set -e

# Write to Obisian
echo -e "- [ ] $text" >>/mnt/f/obsidian/second-brain/Evergreen/todo.md

function insert_to_db {
	dbPath=$1
	sqlite3 "$dbPath" "insert into Tasks (Description, CreatedDate, Priority, Context) values ('$text', '$createdDate', 1, '$context')"
	taskId=$(sqlite3 "$dbPath" "select Id from Tasks where Description = '$text' and CreatedDate = '$createdDate'")
	sqlite3 "$dbPath" "insert into TaskTags (TaskId, TagId) values ($taskId, 5)"
}

testDb="/mnt/x/Stefan/TaskTango/src/config/tasktango.db"
prodDb="$HOME/.tasktango/tasktango.db"

insert_to_db "$prodDb"

if [ -f "$testDb" ]; then
	insert_to_db "$testDb"
fi
