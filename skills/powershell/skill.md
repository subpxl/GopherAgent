# Skill: command-run (PowerShell)

## Description

Executes user requests by breaking them into **single, safe PowerShell commands**, running one command at a time.

This skill is designed for controlled, step-by-step shell execution.

---

## Behavior Rules

1. Interpret the user’s request.
2. Break it into **atomic PowerShell commands** (one action per step).
3. Execute each command sequentially.
4. Show:

   * The command being executed
   * The output
5. Do not combine multiple operations into one command unless necessary.

---

## Execution Format

Each step must follow:

```id="9d9h2a"
Command:
<PowerShell command>

Output:
<result>
```

---

## Safety Constraints

* Do NOT execute destructive commands without explicit confirmation:

  * `Remove-Item`
  * `Clear-Content`
  * `Format-*`
* Avoid system-critical paths unless specified
* No elevation (`Run as Administrator`) unless requested
* Validate paths before execution
* Handle errors gracefully

---

## Common Task Patterns

### List files

```id="z9q1so"
Get-ChildItem
```

### Navigate directory

```id="z1l2sk"
Set-Location "C:\Path\To\Folder"
```

### Create file

```id="g6h2ks"
New-Item -ItemType File -Name "file.txt"
```

### Read file

```id="d3k9sl"
Get-Content "file.txt"
```

### Write to file

```id="c9s8la"
Set-Content "file.txt" "Hello World"
```

---

## Example

### Task

List `.txt` files and read one

### Steps

1.

```id="csc8ak"
Command:
Get-ChildItem -Filter *.txt

Output:
[file list]
```

2.

```id="xv8a2k"
Command:
Get-Content "example.txt"

Output:
[file contents]
```

---

## Error Handling

* If a command fails:

  * Show the error message
  * Suggest a corrected command if possible

---

## Optional Enhancements

* Add dry-run mode (preview commands only)
* Add logging of executed commands
* Add confirmation prompts for risky operations
* Support cross-shell (bash/zsh)

---
