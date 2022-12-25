
# from base64 import decode
# from cgitb import text
# from smtplib import SMTP
# from email.mime.base import MIMEBase
# from email.mime.multipart import MIMEMultipart
# from email import encoders


# sender = 'from@fromdomain.com'
# receivers = ['to@todomain.com']

# mail_content = """From: From Person <from@fromdomain.com>
# To: To Person <to@todomain.com>
# Subject: SMTP e-mail test

# This is a test e-mail message.
# """

# smtp = SMTP()

# message = MIMEMultipart()
# message['Subject'] = 'A test mail sent by Python. It has an attachment.'
# message['From'] = sender
# message['To'] = receivers


# part = MIMEBase('application', "octet-stream")
# part.set_payload(open("hello.txt", "rb").read())
# encoders.encode_base64(part)
# part.add_header('Content-Disposition', 'attachment; filename="text.txt"')

# message.attach(part)

# # text = str(message)

# smtp.connect("smtp.mailtrap.io", 25)
# smtp.login("ca56028e4e358a", "2932a9c05b5fc2")

# decode(message, text)
# smtp.sendmail(sender, receivers, text)
# print("Successfully sent email")

import smtplib
from pathlib import Path
from email.mime.multipart import MIMEMultipart
from email.mime.base import MIMEBase
from email.mime.text import MIMEText
from email.utils import COMMASPACE, formatdate
from email import encoders

sender = 'from@fromdomain.com'
receivers = ['to@todomain.com']
files = ['other/payload_test2.exe']

message = "This is a test e-mail message. For the project.\nThe purpose of the project is detecting mails for phishing."


def send_mail(send_from, send_to, subject, message, files=[],
              server="localhost", port=587, username='', password='',
              use_tls=True):
    """Compose and send email with provided info and attachments.

    Args:
        send_from (str): from name
        send_to (list[str]): to name(s)
        subject (str): message title
        message (str): message body
        files (list[str]): list of file paths to be attached to email
        server (str): mail server host name
        port (int): port number
        username (str): server auth username
        password (str): server auth password
        use_tls (bool): use TLS mode
    """
    msg = MIMEMultipart()
    msg['From'] = send_from
    msg['To'] = COMMASPACE.join(send_to)
    msg['Date'] = formatdate(localtime=True)
    msg['Subject'] = subject

    msg.attach(MIMEText(message))

    # for path in files:
    #     part = MIMEBase('application', "octet-stream")
    #     with open(path, 'rb') as file:
    #         part.set_payload(file.read())
    #     # encoders.encode_base64(part)

    #     part.add_header('Content-Disposition',
    #                     'attachment; filename={}'.format(Path(path).name))
    #     msg.attach(part)

    smtp = smtplib.SMTP(server, port)
    if use_tls:
        smtp.starttls()
    smtp.login(username, password)
    smtp.sendmail(send_from, send_to, msg.as_bytes())
    smtp.quit()


send_mail(sender, receivers, "SMTP e-mail test", message, files,
          "192.168.1.43", 1025, "user", "user", False)
print("Successfully sent email")
