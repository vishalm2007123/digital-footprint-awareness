# digital-footprint-awareness
# OSINT Public Awareness Tool

A lightweight OSINT (Open-Source Intelligence) tool built in **Golang** to help people understand **how much information they expose publicly** on the internet without realizing it.

This project focuses on **awareness and education**, not surveillance, tracking, or intrusion.

---

## 🎯 Purpose

Many people believe that “nothing can be known” from just a username or IP address.  
In reality, attackers, scammers, and analysts often start with **only public data**.

This tool demonstrates:
- What information is visible **without hacking**
- How small pieces of public data can be linked
- Why username reuse and IP exposure matter

It is designed to **inform, not scare**.

---

## 🔍 What This Tool Does

### Username Awareness
- Checks whether a username exists on popular public platforms
- Highlights cross platform username reuse
- Explains why reuse increases linkability and risk

### IP Awareness
- Identifies whether an IP is public or private
- Performs reverse DNS (PTR) lookup
- Explains what infrastructure-level information is exposed

### Exposure Explanation
- Assigns a simple exposure level (Low / Moderate / High)
- Provides clear, human readable warnings
- Explains *why* the information matters

---

## 🚫 What This Tool Does NOT Do

This tool **intentionally avoids** intrusive or unethical behavior.

It does NOT:
- Access private or login protected data
- Use breach databases or leaked credentials
- Perform port scanning or active probing
- Identify real people
- Track users
- Scrape private content
- Use paid or closed APIs

All information shown is **publicly accessible by design**.

---

## 🛠️ How It Works (High Level)

- Uses safe HTTP requests with timeouts
- Relies on public DNS and web responses
- Makes no attempt to bypass protections
- Uses conservative heuristics only

If the tool feels “simple”, that is intentional.

---

## ▶️ Usage

### Requirements
- Go 1.21 or newer

### Run directly
```bash
go run main.go <ip-or-username>
