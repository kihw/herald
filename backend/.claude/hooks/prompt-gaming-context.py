#!/usr/bin/env python3
"""
Herald.lol Gaming Context Injector
Adds gaming platform context to user prompts for better Claude responses
"""

import json
import sys
import re
from datetime import datetime

def inject_gaming_context(prompt: str) -> str:
    """Inject Herald.lol gaming platform context into user prompts"""
    
    # Check if prompt is gaming/analytics related
    gaming_keywords = [
        'analytics', 'gaming', 'league', 'legends', 'tft', 'teamfight', 'tactics',
        'kda', 'cs', 'vision', 'score', 'riot', 'api', 'champion', 'match',
        'player', 'performance', 'metrics', 'dashboard', 'chart'
    ]
    
    is_gaming_related = any(keyword in prompt.lower() for keyword in gaming_keywords)
    
    context_parts = []
    
    # Always add Herald.lol context
    context_parts.append("🎮 Herald.lol Gaming Analytics Platform Context:")
    context_parts.append("- Stack: Go + Gin + React + TypeScript + PostgreSQL + Redis")
    context_parts.append("- Focus: League of Legends & TFT analytics")
    context_parts.append("- Performance targets: <5s analysis, 99.9% uptime, 1M+ concurrent users")
    
    if is_gaming_related:
        context_parts.append("\n📊 Gaming Analytics Focus:")
        context_parts.append("- Core metrics: KDA, CS/min, Vision Score, Damage Share, Gold Efficiency")
        context_parts.append("- Data sources: Riot Games API (rate limited)")
        context_parts.append("- User segments: Amateur → Semi-Pro → Professional")
        context_parts.append("- Architecture: Microservices, Event-driven, Cloud-native")
    
    # Check for specific technical areas
    tech_context = []
    
    if re.search(r'\b(api|endpoint|rest|graphql)\b', prompt, re.IGNORECASE):
        tech_context.append("🌐 API Context: Riot Games API integration with rate limiting")
        
    if re.search(r'\b(database|postgres|redis|migration)\b', prompt, re.IGNORECASE):
        tech_context.append("🗄️ Database Context: PostgreSQL + Redis, gaming data optimization")
        
    if re.search(r'\b(docker|kubernetes|k8s|deploy)\b', prompt, re.IGNORECASE):
        tech_context.append("☸️ Deployment Context: Kubernetes + Docker, auto-scaling for gaming load")
        
    if re.search(r'\b(react|frontend|ui|ux|component)\b', prompt, re.IGNORECASE):
        tech_context.append("🎨 Frontend Context: React + TypeScript + Material-UI, gaming theme")
        
    if re.search(r'\b(performance|optimization|speed|latency)\b', prompt, re.IGNORECASE):
        tech_context.append("⚡ Performance Context: <5s gaming analytics, real-time processing")
        
    if re.search(r'\b(test|testing|coverage|quality)\b', prompt, re.IGNORECASE):
        tech_context.append("🧪 Testing Context: Go tests + Jest, gaming workflow validation")
        
    if re.search(r'\b(security|auth|token|key)\b', prompt, re.IGNORECASE):
        tech_context.append("🔒 Security Context: GDPR compliance, gaming data protection")
    
    if tech_context:
        context_parts.append("\n🔧 Technical Context:")
        context_parts.extend(tech_context)
    
    # Add timestamp for session tracking
    context_parts.append(f"\n⏰ Session time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    
    # Combine context with original prompt
    full_context = "\n".join(context_parts)
    return f"{full_context}\n\n📝 User Request: {prompt}"

def main():
    try:
        input_data = json.load(sys.stdin)
    except json.JSONDecodeError as e:
        print(f"Herald Context Error: Invalid JSON input: {e}", file=sys.stderr)
        sys.exit(1)

    prompt = input_data.get("prompt", "")
    
    if not prompt:
        sys.exit(0)  # No prompt to process
    
    # Inject gaming platform context
    enhanced_prompt = inject_gaming_context(prompt)
    
    # Return enhanced context for Claude
    output = {
        "hookSpecificOutput": {
            "hookEventName": "UserPromptSubmit",
            "additionalContext": enhanced_prompt
        }
    }
    
    print(json.dumps(output))
    sys.exit(0)

if __name__ == "__main__":
    main()