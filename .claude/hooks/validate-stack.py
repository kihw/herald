#!/usr/bin/env python3
"""
Herald.lol Stack Validator
Validates commands against Herald.lol gaming platform best practices
"""

import json
import sys
import re

# Herald.lol specific validation rules
HERALD_VALIDATION_RULES = [
    # Performance and tooling
    (r"\bgrep\b(?!.*\|)", "Use 'rg' (ripgrep) for better performance - gaming analytics require speed"),
    (r"\bfind\s+\S+\s+-name\b", "Use 'rg --files -g pattern' for faster file searching in large codebases"),
    
    # Gaming/Analytics specific
    (r"\bcurl.*riot.*api", "Consider rate limiting for Riot Games API calls - check API key limits"),
    (r"\bgo\s+run\b.*analytics", "Validate gaming metrics before running analytics calculations"),
    (r"\bnpm\s+install\b", "Consider using 'yarn' for consistent gaming platform dependencies"),
    
    # Database and performance
    (r"\bpsql.*-c.*DELETE.*FROM", "Use soft deletes for gaming data to preserve analytics history"),
    (r"\bRECREATE.*DATABASE", "⚠️  Dangerous: Recreating database will lose all gaming analytics data"),
    
    # Security for gaming platform
    (r"\bapi[_-]?key", "⚠️  Potential API key exposure - use environment variables"),
    (r"\briot[_-]?token", "⚠️  Riot API token should be in secure storage"),
    (r"\bpassword.*=.*['\"]", "⚠️  Password in plaintext - use secure credential management"),
    
    # Gaming UX and performance
    (r"\bnpm\s+start", "For gaming platform, use 'npm run dev' with hot reload for better DX"),
    (r"\bdocker.*build.*--no-cache", "Consider using cache for faster builds during gaming feature development"),
]

def validate_gaming_command(command: str) -> list[str]:
    """Validate command against Herald.lol gaming platform requirements"""
    issues = []
    
    for pattern, message in HERALD_VALIDATION_RULES:
        if re.search(pattern, command, re.IGNORECASE):
            issues.append(message)
    
    # Additional gaming-specific validations
    if re.search(r"\btest\b", command, re.IGNORECASE):
        if not re.search(r"(jest|cypress|go\s+test)", command, re.IGNORECASE):
            issues.append("Gaming platform uses Jest (frontend) and Go test (backend) - ensure proper test framework")
    
    if re.search(r"\bbuild\b", command, re.IGNORECASE):
        if re.search(r"production", command, re.IGNORECASE):
            issues.append("✅ Production build detected - ensure gaming analytics are optimized for performance")
    
    return issues

def main():
    try:
        input_data = json.load(sys.stdin)
    except json.JSONDecodeError as e:
        print(f"Herald Hook Error: Invalid JSON input: {e}", file=sys.stderr)
        sys.exit(1)

    tool_name = input_data.get("tool_name", "")
    tool_input = input_data.get("tool_input", {})
    command = tool_input.get("command", "")

    # Only validate Bash commands
    if tool_name != "Bash" or not command:
        sys.exit(0)

    # Validate the command
    issues = validate_gaming_command(command)

    if issues:
        print("🎮 Herald.lol Gaming Platform Validation:", file=sys.stderr)
        for issue in issues:
            print(f"  • {issue}", file=sys.stderr)
        
        # Exit code 2 blocks tool call and shows stderr to Claude
        sys.exit(2)

    # Success - allow command to proceed
    sys.exit(0)

if __name__ == "__main__":
    main()