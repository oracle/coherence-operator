#
# Copyright (c) 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

from typing import List
import jsonpickle
import quart
from coherence import NamedMap, Session, Filters, Processors
from dataclasses import dataclass
from coherence.serialization import proxy
from quart import Quart, request, redirect


@dataclass
class Person:
    id: int
    name: str
    age: int


# ---- init ------------

# the Quart application.  Quart was chosen over Flask due to better
# handling of asyncio which is required to use the Coherence client
# library
app: Quart = Quart(__name__,
                   static_url_path='',
                   static_folder='./')


# the Session with the gRPC proxy
session: Session

people: NamedMap[int, Person]


@app.before_serving
async def init():

    # initialize the session using the default localhost:1408 or the value of COHERENCE_SERVER_ADDRESS
    global session
    session = await Session.create()

    global people
    people = await session.get_map('people')

# ----- routes --------------------------------------------------------------

# Get all people
@app.route('/api/people', methods=['GET'])
async def get_people():
    people_list: List[People] = []
    async for person in await people.values():
        people_list.append(person)

    return quart.Response(jsonpickle.encode(people_list, unpicklable=False), mimetype="application/json")

# Create a person with JSON as body
@app.route('/api/people', methods=['POST'])
async def create_person():
    data = await request.get_json(force=True)
    name: str = data['name']
    id: int = data['id']
    age: int = data['age']
    person: Person = Person(id, name, age)
    await people.put(person.id, person)

    return quart.Response(
        jsonpickle.encode(person, unpicklable=False),
        status=201,
        mimetype='application/json'
    )

# Get a single person
@app.route('/api/people/<id>', methods=['GET'])
async def get_person(id: str):
    existing: Person = await people.get(int(id))
    if existing == None:
        return "", 404

    return jsonpickle.encode(existing, unpicklable=False), 200

# Delete a person
@app.route('/api/people/<id>', methods=['DELETE'])
async def delete_person(id: str):
    """
    This route will delete the person with the given id.

    :param id: the id of the person to delete
    """
    existing: Person = await people.remove(int(id))
    return "", 404 if existing is None else 200

# ----- main ----------------------------------------------------------------

if __name__ == '__main__':
    # run the application on port 8080
    app.run(host='0.0.0.0', port=8080)