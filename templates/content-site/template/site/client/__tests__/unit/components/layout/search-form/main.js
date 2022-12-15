import * as directives from 'sham-ui-directives';
import LayoutSearchFormMain  from '../../../../../src/components/layout/search-form/main.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( LayoutSearchFormMain, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
