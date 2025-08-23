import sys

def count_numbers(filename):
    with open(filename, "r") as f:
        lines = [line.strip() for line in f if line.strip().isdigit()]

    i = 0
    while i < len(lines):
        num = int(lines[i])

        # Special stop condition
        if num in (10, 100):
            print(f"code {num} was triggered flags set to :")
            if i + 1 < len(lines):
                print(lines[i + 1])  # print next number once
            break

        # Count consecutive occurrences
        count = 1
        while i + 1 < len(lines) and lines[i + 1] == lines[i]:
            count += 1
            i += 1

        print(f"{num} : {count}")
        i += 1




if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python script.py <filename>")
    else:
        count_numbers(sys.argv[1])
