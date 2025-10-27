const express = require('express');
const morgan = require('morgan');

const app = express();
const PORT = process.env.PORT || 3000;

app.use(morgan('dev'));
app.use(express.json());

app.get('/healthz', (req, res) => {
  res.json({ status: 'ok' });
});

app.get('/api/query', (req, res) => {
  const { name = 'unknown', age = 'unknown' } = req.query;
  res.json({
    message: 'Received query parameters from gateway',
    name,
    age,
    timestamp: new Date().toISOString()
  });
});

app.get('/api/user/:id', (req, res) => {
  res.json({
    message: 'Fetched user data',
    id: req.params.id,
    user: {
      id: req.params.id,
      name: 'Gateway User',
      joinedAt: '2024-01-01'
    }
  });
});

app.post('/api/create', (req, res) => {
  res.status(201).json({
    message: 'Resource created',
    payload: req.body ?? {},
    requestId: req.header('x-request-id') || 'generated-by-example-backend-service'
  });
});

app.delete('/api/delete/:id', (req, res) => {
  res.json({
    message: 'Resource deleted',
    id: req.params.id
  });
});

// The deploy.yaml path currently has a trailing space, so we accept both forms.
app.put('/api/update/:id', (req, res) => {
  res.json({
    message: 'Resource updated',
    id: req.params.id,
    payload: req.body ?? {}
  });
});

app.put('/api/update/:id ', (req, res) => {
  res.json({
    message: 'Resource updated (trimmed path)',
    id: req.params.id,
    payload: req.body ?? {}
  });
});

app.listen(PORT, () => {
  // eslint-disable-next-line no-console
  console.log(`Example service listening on port ${PORT}`);
});
