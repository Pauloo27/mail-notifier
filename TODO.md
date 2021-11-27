# To-Do

- [ ] daemon:
  - [ ] encrypt token/user/password
  - [ ] socket:
    - [ ] handle new connections
    - [ ] handle requests (see #format):
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
  - [.] gmail:
    - [ ] parse mail to list with regex?
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
- `fetch_unread_messages [id]`: return unread messages for a inbox.
- `fetch_message [inbox_id] [id]`: return the data for a single message.
- `mark_as_message [inbox_id] [id]`: return the data for a single message.

