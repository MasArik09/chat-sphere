# UI/UX Flow

# ChatSphere V1

Version: 1.0

Status: Approved

---

# 1. Design Principles

ChatSphere UI must be:

* Clean
* Modern
* Responsive
* Familiar
* Fast

Design inspiration:

* WhatsApp Web
* Telegram Web
* Discord Direct Messages

Avoid:

* Complex dashboards
* Excessive animations
* Fancy effects

Focus:

Messaging first.

---

# 2. Application Layout

Desktop Layout

+------------------------------------------------------+
| Sidebar | Conversation List | Chat Window           |
+------------------------------------------------------+

---

Mobile Layout

+-----------------------+
| Conversation List     |
+-----------------------+

↓

Tap Conversation

↓

+-----------------------+
| Chat Window           |
+-----------------------+

---

# 3. Authentication Pages

## Login Page

Purpose:

Allow existing users to login.

Layout:

+--------------------------------+
|          ChatSphere            |
|                                |
| Email                          |
| [******************]           |
|                                |
| Password                       |
| [******************]           |
|                                |
| [ Login ]                      |
|                                |
| Register Link                  |
+--------------------------------+

---

Actions:

Login

Navigate to Register

---

## Register Page

Layout:

+--------------------------------+
|          ChatSphere            |
|                                |
| Name                           |
| [******************]           |
|                                |
| Email                          |
| [******************]           |
|                                |
| Password                       |
| [__________________]           |
|                                |
| [ Register ]                   |
+--------------------------------+

---

# 4. Main Application Layout

Desktop

+-------------------------------------------------------------+
| Sidebar | Conversations | Chat Window                      |
+-------------------------------------------------------------+

Sidebar Width:

240px

Conversation List Width:

320px

Chat Window:

Remaining Width

---

# 5. Sidebar

Purpose:

Global navigation.

Contains:

* Logo
* Current User
* Online Status
* Logout

Layout:

+--------------------+
| ChatSphere         |
|                    |
| John Doe           |
| ● Online           |
|                    |
| Logout             |
+--------------------+

---

# 6. Conversation List

Purpose:

Display user conversations.

Layout:

+--------------------------------+
| Search                         |
| [______________]               |
+--------------------------------+

Conversation Card

+--------------------------------+
| Jane Smith                     |
| Last message preview...        |
| 12:45 PM                       |
+--------------------------------+

---

Display:

* Username
* Last Message
* Timestamp

Ordered:

Most recent first

---

# 7. User Directory

Purpose:

Find users.

Layout:

+--------------------------------+
| Search Users                   |
| [______________]               |
+--------------------------------+

User Card

+--------------------------------+
| Jane Smith      ● Online       |
| [ Start Chat ]                 |
+--------------------------------+

---

Actions:

Start Chat

---

# 8. Chat Window

Purpose:

Display conversation messages.

Layout:

+------------------------------------------------+
| Jane Smith                     ● Online        |
+------------------------------------------------+

|                                                |
| Message History                                |
|                                                |
|  Hello                                         |
|                                  Hi there      |
|                                                |
|                                                |
+------------------------------------------------+

| [Type message...]      [Send]                  |
+------------------------------------------------+

---

Sections:

Header

Messages

Composer

---

# 9. Chat Header

Display:

* User Name
* Online Status

Examples:

Jane Smith

● Online

or

Last Seen 5 minutes ago

---

# 10. Message Bubbles

Current User

Align Right

Style:

Primary Color

---

Other User

Align Left

Style:

Neutral Color

---

Display:

* Content
* Time

Example

Hello!

12:30 PM

---

# 11. Message Composer

Layout:

+---------------------------------------+
| Type message...          [Send]       |
+---------------------------------------+

Rules:

Empty messages disabled.

Send on:

Button Click

Enter Key

---

# 12. Presence Indicators

Online

Green Dot

● Online

---

Offline

Gray Dot

● Offline

---

Last Seen

Last seen 10 minutes ago

---

# 13. Empty States

## No Conversations

+--------------------------------+
| No Conversations Yet           |
|                                |
| Start chatting with someone.   |
|                                |
| [ Find Users ]                 |
+--------------------------------+

---

## No Messages

+--------------------------------+
| No Messages Yet                |
|                                |
| Say hello 👋                   |
+--------------------------------+

---

## No Search Results

+--------------------------------+
| No Users Found                 |
+--------------------------------+

---

# 14. Loading States

Conversation List

Loading...

---

Messages

Loading messages...

---

User Directory

Loading users...

---

Use skeleton loaders where possible.

---

# 15. Error States

Network Error

+--------------------------------+
| Unable to load data            |
|                                |
| [ Retry ]                      |
+--------------------------------+

---

Connection Lost

+--------------------------------+
| Reconnecting...                |
+--------------------------------+

---

# 16. Responsive Behavior

Desktop

Sidebar visible

Conversation list visible

Chat visible

---

Tablet

Smaller widths

Layout unchanged

---

Mobile

Conversation list first

Chat opens separately

Back button available

---

# 17. Color System

Primary

Indigo

---

Success

Green

---

Danger

Red

---

Neutral

Slate / Gray

---

Presence

Green Dot

Gray Dot

---

# 18. Accessibility

Must support:

* Keyboard navigation
* Focus states
* Visible labels
* Color contrast

---

# 19. Page Flow

Login

↓

Dashboard

↓

Conversation List

↓

Open Conversation

↓

Realtime Messaging

---

User Search

↓

Start Chat

↓

Conversation Created

↓

Realtime Messaging

---

# 20. Definition of Done

UI accepted when:

✓ Login page complete

✓ Register page complete

✓ User directory complete

✓ Conversation list complete

✓ Chat window complete

✓ Presence indicators complete

✓ Empty states complete

✓ Mobile layout defined
