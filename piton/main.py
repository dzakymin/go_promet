import random
import math

def is_prime(number):
    if number < 2:
        return False
    for i in range(2, number // 2 + 1):
        if number % i == 0:
            return False
    return True

def generate_prime(min_number,max_number):
    prime = random.randint(min_number,max_number)
    while not is_prime(prime):
        prime = random.randint(min_number,max_number)
    return prime

def mod_inverse(e, phi):
    for d in range(3, phi):
        if (e * d) % phi == 1:
            return d
    raise ValueError("Error value")


p,q = generate_prime(1000,3000), generate_prime(1000,3000)

while p == q:
    q = generate_prime(1000,2000)

N = p * q
phi_n = (p-1) * (q-1)

e = random.randint(3, phi_n-1)

while math.gcd(phi_n, e) != 1:
    e = random.randint(3, phi_n-1)

d = mod_inverse(e, phi_n)

massage = "hello warudo"

massage_encode = [ord(ch) for ch in massage]
cipertext = [pow(ch, e, N) for ch in massage_encode]
print(cipertext)


real_text = [pow(ch, d, N) for ch in cipertext]
real = [ord(ch) for ch in real_text]

print(real)

