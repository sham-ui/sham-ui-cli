import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
    setup();
} );

it( 'fail edit name', async() => {
    expect.assertions( 1 );

    axios.useDefaultMocks();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start();
    app.click( '.panel.settings p:nth-of-type(1) .icon-pencil' );

    const formData = {
        name: ''
    };
    axios
        .use( 'put', '/members/name', {
            'Status': 'Bad Name',
            'Messages': [ 'Name must have more than 0 characters.' ]
        }, 400 );

    app.form.fill( 'name', formData.name );
    await app.form.submit();

    app.checkBody();
} );
