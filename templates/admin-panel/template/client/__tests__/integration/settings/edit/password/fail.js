import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'fail edit password', async() => {
    expect.assertions( 1 );

    axios.useDefaultMocks();
    const DI = setup();

    history.pushState( {}, '', 'http://client.example.com/settings/' );
    await app.start( DI );
    app.click( '.panel.settings p:nth-of-type(3) .icon-pencil' );

    const formData = {
        newPass1: 'test',
        newPass2: 'test1'
    };
    axios
        .use( 'put', '/members/password', {
            'Status': 'Bad Password',
            'Messages': [ 'Passwords don\'t match.' ]
        }, 400 );

    app.form.fill( 'pass1', formData.newPass1 );
    app.form.fill( 'pass2', formData.newPass2 );
    await app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );

    await app.waitRendering();

    app.checkBody();
} );
