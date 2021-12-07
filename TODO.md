# To-Do

- [ ] client:
  - [ ] handle server not running/connection lost
- [ ] daemon:
  - [ ] encrypt token/user/password
  - [ ] use generics (go 1.18) for a serialization func
  - [ ] socket client:
    - [ ] accept "force update" command
  - [ ] socket server:
    - [ ] trigger "force update"
    - [ ] auto refresh cache
- [ ] cli:
  - [ ] polybar integration
- [ ] gui:
  - [ ] settings page:
    - [ ] new account btn
    - [ ] account listing (with delete button)
  - [O] inbox page:
    - [ ] preview message (webview? nope. cant disable js/image load)
- [O] core:
  - [o] gmail:
    - [ ] batch-get messages?
  - [o] imap:
    - [ ] batch-get messages?

# Format

The deamon will read data until it receives a `\n`. The data received will be
interpreted as a command. The available commands are:

- `list_inboxes`: return the list of inboxes (id, address).
- `fetch_message [inbox_id] [id]`: return the data for a single message.
- `fetch_unread_messages_in [id]`: return unread messages for a inboxes.
- `fetch_all_unread_messages [id]`: return unread messages for all inboxes.
- `mark_message_as_read [inbox_id] [id]`: mark a message as read.
- `refresh_inbox [inbox_id]`: return the data for a single message.
