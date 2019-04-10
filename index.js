require("request");
const request = require("request-promise");
const { parse } = require("url");

async function createAndDownloadImage() {
  const data = {
    html: `<div class='box'>${Date.now()}</div>`,
    css: ".box { border: 4px solid #03B875; padding: 20px; font-family: 'Roboto'; }",
    google_fonts: "Roboto"
  };

  try {
    const createImage = await request
      .post({
        url: "https://hcti.io/v1/image",
        form: data,
        time: true,
        resolveWithFullResponse: true
      })
      .auth(process.env.API_ID, process.env.API_KEY);

    const { url } = JSON.parse(createImage.body);

    const downloadImage = await request.get({
      url,
      time: true,
      resolveWithFullResponse: true
    });

    return `<pingdom_http_custom_check>
    <status>OK</status>
    <response_time>${createImage.elapsedTime + downloadImage.elapsedTime}</response_time>
</pingdom_http_custom_check>`;
  } catch (err) {
    console.log(err);

    return `<pingdom_http_custom_check>
    <status>DOWN</status>
</pingdom_http_custom_check>`;
  }
}

module.exports = async (req, res) => {
  const { query = {} } = parse(req.url, true);

  if (query["password"] !== process.env.UPTIME_PASSWORD) {
    res.statusCode = 401;
    res.end("Unauthorized");
  } else {
    const xml = await createAndDownloadImage();

    res.setHeader("Content-Type", "text/xml");
    res.end(xml);
  }
};
