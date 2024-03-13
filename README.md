# ImgEncryptor

Image Encryptor is a command-line tool written in Go for securing image files through encryption. It provides functionalities for both encryption and decryption, ensuring the privacy and integrity of visual data from unauthorized access and tampering.
Features:

    Encryption: Encrypt image files using AES encryption algorithm in CTR mode.
    Decryption: Decrypt encrypted image files using the provided encryption key.
    Key Management: Automatically generates and stores the encryption key in a file named key.txt.
    User-friendly Interface: Interactive CLI interface for easy usage.

Usage:

    Generate a Key: If you haven't generated a key yet, run the program. It will automatically generate a new key and save it to key.txt.
    Encrypt an Image: Run the program in encrypt mode, providing the path to the input image file and the path to save the encrypted image file.
    Decrypt an Image: Run the program in decrypt mode, providing the path to the input encrypted image file and the path to save the decrypted image file.
