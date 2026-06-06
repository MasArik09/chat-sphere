# ChatSphere V1 - Screenshots & Demo Guide

This document outlines the recommended procedures for capturing visual media (screenshots and GIFs) to showcase ChatSphere V1 on a GitHub repository README.

---

## 1. Screenshots Checklist

To provide a comprehensive visual overview of the application, capture the following screenshots (in high-resolution PNG format, using 16:9 aspect ratios):

1. **User Authentication Interface**:
   - **Target**: Login or Register screen displaying clean input validation and modern form styling.
   - **Details**: Ensure placeholder values look realistic (e.g. `user@example.com`).
2. **Conversations Dashboard**:
   - **Target**: The main user chat screen showing the active chat threads list, search bar, unread indicators, and user online badges.
   - **Details**: Mock a rich chat list showing multiple user avatars, previews of last messages, and green/gray presence indicators.
3. **Active Chat Window**:
   - **Target**: A split screen interface showing a conversation between two users with different message alignments, read receipts, and timestamps.
   - **Details**: Show clear sender/recipient styling (e.g., right-aligned user messages, left-aligned partner messages).
4. **Mobile Responsive View**:
   - **Target**: The conversations view adjusted on a mobile device emulator (e.g., iPhone 15 Pro dimensions).
   - **Details**: Highlight the responsive sidebar folding behaviors and mobile-optimized chat windows.

---

## 2. Recommended Demo Flow

To demonstrate the real-time functionality of ChatSphere, open two browser windows side-by-side (e.g., Window A logged in as *Alice*, Window B logged in as *Bob*):

```text
+----------------------------+  +----------------------------+
|        Window A: Alice      |  |         Window B: Bob      |
|  Logged in, looking at list |  |  Logged in, looking at list |
+----------------------------+  +----------------------------+
```

### Steps:
1. **Presence Verification**: Note the active green presence indicator next to Bob's avatar inside Alice's window, and vice-versa.
2. **Typing Indicators**: In Window A, Alice starts typing. Bob immediately sees "Alice is typing..." indicator inside Window B. Alice stops typing, and the indicator immediately disappears.
3. **Real-time Delivery**: Alice sends a message: *"Hi Bob!"*. The message instantly appears in Bob's Window B without page refresh.
4. **Read Receipt Sync**: Bob clicks on the chat to view the message. Alice's Window A instantly updates to display a visual read receipt matching Bob's user ID.
5. **Presence Toggle**: Close Window B (Bob logs out or disconnects). Alice's Window A instantly updates Bob's indicator to offline (gray badge) and lists his last-seen timestamp.

---

## 3. Recommended GIF Recording Sequence

For dynamic presentation, capture short (~5-10s) high-quality looping GIFs for these key workflows:

### GIF 1: Real-Time Typing & Message Synchronization
- **Action**: Alice typing -> Bob seeing typing state -> Alice sending -> Bob instantly receiving the text message.
- **Duration**: ~6 seconds.

### GIF 2: Reactive Online/Offline Presence
- **Action**: Closing Bob's window -> Alice's presence list immediately updating Bob's status badge from green (online) to gray (offline).
- **Duration**: ~4 seconds.

### GIF 3: Real-Time Read Receipts Update
- **Action**: Bob opening a conversation -> Alice's window immediately displaying the "Read by Bob" visual indicator under the message.
- **Duration**: ~4 seconds.

---

## 4. Suggested GitHub README Media Layout

To present screenshots and GIFs cleanly without cluttering the README, utilize markdown table structures and HTML image wrapping:

```markdown
### 📱 Application in Action

<table width="100%">
  <tr>
    <td width="50%" align="center">
      <b>Real-Time Chat & Typing Synchronization</b><br/>
      <img src="docs/media/typing_demo.gif" alt="Typing Sync GIF" width="100%"/>
    </td>
    <td width="50%" align="center">
      <b>Online Presence Tracking</b><br/>
      <img src="docs/media/presence_demo.gif" alt="Presence Sync GIF" width="100%"/>
    </td>
  </tr>
</table>

### 📸 UI Screenshots

<details>
  <summary><b>View Dashboard & Mobile Layout Screenshots</b></summary>
  <br/>
  
  #### Dashboard View
  ![Dashboard](docs/media/dashboard_view.png)
  
  #### Mobile Layout
  ![Mobile](docs/media/mobile_view.png)
</details>
```
