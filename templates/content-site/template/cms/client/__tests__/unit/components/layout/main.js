import { createDI } from 'sham-ui';
import LayoutMain  from '../../../../src/components/layout/main.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();

    DI.bind( 'router', {
        generate: jest.fn().mockReturnValueOnce( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/'
        }
    } );

    const meta = renderer( LayoutMain, { DI } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
