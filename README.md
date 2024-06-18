# TaskTango

Dance around your tasks

Short description: This is an opionated task management tool, inspired by taskwarrior.
The main purpouse is to have a centralized place for all tasks with some basic features
to manage them and whit frictionless capture.

## Objectives

### 1. Manage tasks - done

It should be able to quickly retain tasks that I have to work on,
These tasks can be ordered by priority
Tasks can be marked in progress, can be soft deleted once finished.
Tasks should have a project assigned, and tags.

For the managing part, you should have a view of all your tasks with pagination
Multiple options of sorting and filtering
Ability to do CRUD on tasks.

There should be a quick capture flow, to get task description
There should be a review feature to better describe tasks with additional metadata
for example at the end of the day after capturing multiple tasks during the day.
There should be some statistics (time spent on task)

### 2. Allow input of tasks from phone

Very little support on any phone with internet
Just to send a new task to the system for quick capture while away from PC
That's it.

### 3. Sync and backup

Allow storing data in the cloud, encrypted.
You should be able to back it up periodically to avoid losing it.

### 4. Integrations

Possible integrations, integrate with Outlook calendar to book slots
when a task will be in progress

## Dependencies

Bubbletea
Bubble-table
Bubbles
Huh
Lipgloss
SQLite3
