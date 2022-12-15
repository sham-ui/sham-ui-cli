import * as directives from 'sham-ui-directives';
import { createDI } from 'sham-ui';
import RoutesTagPage  from '../../../../../src/components/routes/tag/page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );

    DI.bind( 'store', {
        articles: jest.fn().mockReturnValueOnce(
            Promise.resolve( { meta: {} } )
        )
    } );

    DI.bind( 'router', {
        storage: {
            params: {
                page: 2
            }
        },
        generate: () => '/'
    } );

    const meta = renderer( RoutesTagPage, {}, {
        DI,
        directives
    } );

    await new Promise( resolve => setImmediate( resolve ) );

    expect( meta.toJSON() ).toMatchSnapshot();
} );
