/*
 * Copyright (c) 2025 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

const coh = require('@oracle/coherence')

const Session = coh.Session

const express = require('express');
const port = process.env.PORT || 8080
const api = express();

api.use(express.json()); // to parse JSON request bodies

// setup session to Coherence
const session = new Session()
const people = session.getCache('people')

// ----- REST API -----------------------------------------------------------

/**
 * Returns all people.
 */
api.get('/api/people', (req, res, next) => {
    const toSend = []
    people.values()
        .then(async values => {
            // copy values to array to be sent via express
            for await (let value of values) {
                toSend.push(value)
            }
            res.send(toSend)
        })
        .catch(err => next(err))
})


/**
 * Create a person.
 */
api.post('/api/people', (req, res, next) => {
    const id = req.body.id
    const person = {
        id: req.body.id,
        name: req.body.name,
        age: req.body.age
    }

    people.set(id, person)
        .then(() => {
            res.send(JSON.stringify(person))
        })
        .catch(err => next(err))
})


/**
 * Get a single person.
 */
api.get('/api/people/:id', (req, res, next) => {
    const id = Number(req.params.id);
    people.get(id)
        .then(person => {
            if (person) {
                res.status(200).json(person);
            } else {
                res.sendStatus(404);
            }
        })
        .catch(err => next(err));
});

/**
 * Delete a person.
 */
api.delete('/api/people/:id', (req, res, next) => {
    const id = Number(req.params.id);
    people.delete(id)
        .then(oldValue => {
            res.sendStatus(oldValue ? 200 : 404)
        })
        .catch(err => next(err))
})

api.listen(port, () => console.log(`Listening on port ${port}`))
