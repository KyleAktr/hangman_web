function sendMail() {
  let parms = {
    message: document.getElementById("message").value,
  };

  emailjs
    .send("service_f8b91il", "template_zb29n8b", parms)
    .then(alert("email envoyé !"));
}
