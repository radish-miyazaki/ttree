# ttree Demo
Output demo.gif

Set Shell zsh
Set Width 1200
Set Height 600
Set FontSize 16
Set Padding 20
Set TypingSpeed 60ms

# Launch ttree
Type@0ms "./ttree"
Enter
Sleep 2s

# Type root node
Type "Project"
Sleep 500ms

# Add child: src
Enter
Sleep 300ms
Type "src"
Sleep 300ms
Tab
Sleep 500ms

# Add child of src: main.go
Enter
Sleep 300ms
Type "main.go"
Sleep 300ms
Tab
Sleep 500ms

# Add sibling: utils.go
Enter
Sleep 300ms
Type "utils.go"
Sleep 500ms

# Go back up to Project level, add docs
Enter
Sleep 300ms
Shift+Tab
Sleep 300ms
Type "docs"
Sleep 300ms
Tab
Sleep 500ms

# Add child of docs: README.md
Enter
Sleep 300ms
Type "README.md"
Sleep 300ms
Tab
Sleep 500ms

# Go back, add tests
Enter
Sleep 300ms
Shift+Tab
Sleep 300ms
Type "tests"
Sleep 500ms

# Add child of tests: main_test.go
Enter
Sleep 300ms
Type "main_test.go"
Sleep 300ms
Tab
Sleep 1s

# Navigate up to show cursor movement
Up
Sleep 300ms
Up
Sleep 300ms
Up
Sleep 300ms
Up
Sleep 300ms
Up
Sleep 300ms
Up
Sleep 300ms
Up
Sleep 1s

# Navigate back down
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 1s

# Copy tree to clipboard
Ctrl+C
Sleep 2s

# Quit
Ctrl+Q
Sleep 1s
