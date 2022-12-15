import * as directives from 'sham-ui-directives';
import PaginationPages  from '../../../../src/components/pagination/pages.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( PaginationPages, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
