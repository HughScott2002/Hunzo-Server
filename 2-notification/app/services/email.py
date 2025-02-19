from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
import smtplib
from jinja2 import Template
from app.core.config import settings

def render_template(template_name: str, context: dict) -> str:
   with open(f"templates/{template_name}.html") as f:
       template = Template(f.read())
   return template.render(context)

def send_email(to: str, template_name: str, context: dict):
   html_content = render_template(template_name, context)
   
   message = MIMEMultipart("alternative")
   message["Subject"] = context.get("subject", "Notification")
   message["From"] = settings.email_from
   message["To"] = to
   
   html_part = MIMEText(html_content, "html")
   message.attach(html_part)

   with smtplib.SMTP(settings.smtp_server, settings.smtp_port) as server:
       server.starttls()
       server.login(settings.email_from, settings.smtp_password)
       server.send_message(message)