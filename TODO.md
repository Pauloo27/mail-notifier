# To-Do

- [.] daemon:
  - [ ] encrypt token/user/password
  - [ ] use generics (go 1.18) for a serialization func
  - [o] socket client:
    - [X] command sender
    - [ ] accept "force update" command
  - [O] socket server:
    - [o] trigger "force update"
    - [X] handle new connections
    - [X] handle requests (see #format):
    - [ ] auto refresh cache
    - [X] commands:
      - [X] list inboxes
      - [X] fetch message
      - [X] fetch unread messages (per inbox/all)
- [ ] cli:
  - [ ] polybar integration
- [ ] gui:
  - [X] home page:
    - [X] list inboxes 
    - [X] open inbox (launch browser)
    - [X] loading spinner
  - [ ] settings page:
    - [ ] new account btn
    - [ ] account listing (with delete button)
  - [O] inbox page:
    - [X] list messages 
    - [X] mark as read btn
    - [X] loading spinner
    - [X] async mark as read
    - [X] add ... to long subjects
    - [ ] preview message (webview? nope. cant disable js/image load)
    - [X] open message (launch browser)
- [O] core:
  - [X] create provider interface
  - [X] mark as read interface
  - [o] gmail:
    - [X] parse mail to list with regex?
    - [X] mark as read
    - [ ] batch-get messages?
  - [o] imap:
    - [X] load info from file
    - [X] mark as read impl
    - [ ] batch-get messages?

# Format

The deamon will read data until it receives a `\n`. The data received will be
interpreted as a command. The available commands are:

- `list_inboxes`: return the list of inboxes (id, address).
- `fetch_message [inbox_id] [id]`: return the data for a single message.
- `fetch_unread_messages_in [id]`: return unread messages for a inboxes.
- `fetch_all_unread_messages [id]`: return unread messages for all inboxes.
- `mark_as_message [inbox_id] [id]`: return the data for a single message.
- `refresh_inbox [inbox_id]`: return the data for a single message.
