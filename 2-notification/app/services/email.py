from sendgrid import SendGridAPIClient
from jinja2 import Template
from app.core.config import settings

def render_template(template_name: str, context: dict) -> str:
    with open(f"templates/{template_name}.html") as f:
        template = Template(f.read())
    return template.render(context)

def send_email(to: str, template_name: str, context: dict):
    html_content = render_template(template_name, context)
    message = {
        "to": [{"email": to}],
        "from": {"email": settings.email_from},
        "subject": context.get("subject", "Notification"),
        "html": html_content
    }
    sg = SendGridAPIClient(settings.sendgrid_api_key)
    sg.send(message)