import random
import string

C = 50
for i in range(C + (C // 3)):
    print("".join(random.choice(string.ascii_uppercase) for _ in range(C)))
