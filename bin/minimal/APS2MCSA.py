import re

def convert_ports(ip_address, input_file="raw.txt", output_file="targets.txt"):
    try:
        with open(input_file, "r", encoding="utf-8") as f:
            content = f.read()
    except FileNotFoundError:
        print(f"Ошибка: Создай файл {input_file} и закинь туда вывод сканера!")
        return

    # Регулярка ищет слово Port, а затем забирает только цифры
    ports = re.findall(r"Port\s+(\d+)", content)

    if not ports:
        print("Порты не найдены. Проверь формат текста в файле.")
        return

    # Собираем строки вида IP:PORT
    result_lines = [f"{ip_address}:{port}" for port in ports]

    with open(output_file, "w", encoding="utf-8") as f:
        f.write("\n".join(result_lines) + "\n")

    print(f"Готово! Обработано портов: {len(ports)}")
    print(f"Результат сохранен в {output_file}")

if __name__ == "__main__":
    target_ip = input("Target IP: ")
    convert_ports(target_ip)
    input()