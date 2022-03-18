const http = require('http');
const yt = require('youtube-search-without-api-key');

const hostname = '127.0.0.1';
const port = 3000;

const server = http.createServer( async (req, res) => {
  res.statusCode = 200;
  res.setHeader('Content-Type', 'text/html; charset=utf-8');
  var query = require('url').parse(req.url,true).query;
  const videos = await yt.search(query.q);
  res.write(JSON.stringify(videos));
});

server.listen(port, hostname, () => {
  console.log(`Server running at http://${hostname}:${port}/`);
});

