# timer-timer
a terminal-based multi-timer app written in Go using tview. Features hidden countdowns, ASCII 'falling-sand' visualizations, and modal notifications - built to reduce time-checking anxiety and make waiting more fun.
imer/timer

A terminal-based multi-timer application written in Go using the tview TUI framework.

This project started as a small experiment while learning Go: I wanted a timer that wouldn’t constantly show me how much time I had left. Instead of stressful countdown numbers, this app uses an ASCII “falling sand” visualization to represent remaining time. Longer timers fall slowly, shorter timers fall quickly, and each timer has its own uniquely colored sand animation.

If you truly want to see the remaining time, you can — but only by navigating through deliberate steps so you never accidentally glimpse it.

## Features

⏱ Multiple simultaneous timers
Create and manage as many timers as you want, each with its own label.

🌈 ASCII “falling sand” time representation
A dynamic, color-based animation showing how much time is left without revealing exact numbers.

🙈 Intentional hidden countdown
Remaining time is only visible if you specifically navigate to it.

🔔 Modal completion notifications
When a timer finishes, a blocking modal appears so you can’t miss it.

⏸ Pause / resume / stop timers
Basic timer controls for flexibility.

💾 No persistence by design
Timers do not persist between sessions, keeping the app simple and distraction-free.

### Why I Built This

This was my first “real” project while learning Go.
I wanted something small, useful, and fun that solved a personal annoyance: I checked countdown timers too often. The falling-sand idea was a playful alternative — a way to feel time passing without obsessing over numbers.

### Roadmap / Future Ideas

Active development may slow down, but here are some ideas I might revisit:

Additional color themes

Adjustable sand density / speed

Improved keybindings

Configurable defaults

General refactoring and cleanup
