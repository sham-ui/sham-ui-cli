import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import RoutesSignupPage  from '../../../../../src/components/routes/signup/page.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';
import { createDI } from 'sham-ui';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );
    const meta = renderer( RoutesSignupPage, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );


it( 'display errors', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    DI.bind( 'router', {
        DI,
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );

    const signUpMock = jest.fn();
    DI.bind( 'store', {
        signUp: signUpMock.mockReturnValueOnce( Promise.reject( {} ) )
    } );

    const meta = renderer( RoutesSignupPage, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );

    const formData = {
        name: 'admin',
        email: 'admin@gmail.com',
        password: 'passw0rd'
    };
    const { component } = meta;
    component.container.querySelector( '[name="name"]' ).value = formData.name;
    component.container.querySelector( '[name="email"]' ).value = formData.email;
    component.container.querySelector( '[name="password"]' ).value = formData.password;
    component.container.querySelector( '[name="password2"]' ).value = formData.password;
    component.container.querySelector( '[type="submit"]' ).click();

    await new Promise( resolve => setImmediate( resolve ) );

    expect( signUpMock.mock.calls ).toHaveLength( 1 );
    expect( signUpMock.mock.calls[ 0 ] ).toHaveLength( 1 );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
