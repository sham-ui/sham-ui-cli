import PageLoadingIndicator  from '../../../src/components/PageLoadingIndicator.sht';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( PageLoadingIndicator, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
