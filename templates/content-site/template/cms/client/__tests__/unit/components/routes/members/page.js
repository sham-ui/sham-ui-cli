import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import RoutesMembersPage  from '../../../../../src/components/routes/members/page.sfc';
import renderer from 'sham-ui-test-helpers';
import { storage } from '../../../../../src/storages/session';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;

    const getMock = jest.fn();
    DI.bind( 'store', {
        api: {
            request: getMock.mockReturnValueOnce(
                Promise.resolve( { meta: {} } )
            )
        }
    } );

    const meta = renderer( RoutesMembersPage, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'store', {
        api: {
            request: jest.fn().mockReturnValueOnce(
                Promise.reject( {} )
            )
        }
    } );
    const meta = renderer( RoutesMembersPage, {
        DI,
        directives
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
