import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import IdleTimer from 'react-idle-timer';
import { isProduction } from 'shared/constants';
import Alert from 'shared/Alert';

const fifteenMinutesinMilliseconds = 900000;
const oneMinuteInMilliseconds = 60000;
export class LogoutOnInactivity extends React.Component {
  state = {
    isIdle: false,
    showLoggedOutAlert: false,
  };
  componentDidUpdate(prevProps) {
    if (!this.props.isLoggedIn && prevProps.isLoggedIn) {
      this.setState({ showLoggedOutAlert: true });
    }
  }
  onActive = () => {
    this.setState({ isIdle: false });
  };
  onIdle = () => {
    this.setState({ isIdle: true });
    setInterval(
      () => (window.location = '/auth/logout'),
      oneMinuteInMilliseconds,
    );
  };
  render() {
    const props = this.props;
    return (
      <React.Fragment>
        {isProduction &&
          props.isLoggedIn && (
            <IdleTimer
              ref="idleTimer"
              element={document}
              activeAction={this.onActive}
              idleAction={this.onIdle}
              timeout={fifteenMinutesinMilliseconds}
            >
              {this.state.isIdle && (
                <div className="usa-grid">
                  <Alert type="warning" heading="Inactive user">
                    You have been inactive and will be logged out shortly.
                  </Alert>
                </div>
              )}
            </IdleTimer>
          )}

        {this.state.showLoggedOutAlert && (
          <div className="usa-grid">
            <Alert type="error" heading="Logged out">
              You have been logged out due to inactivity.
            </Alert>
          </div>
        )}
      </React.Fragment>
    );
  }
}
LogoutOnInactivity.defaultProps = {
  idleTimeout: fifteenMinutesinMilliseconds,
};
LogoutOnInactivity.propTypes = {
  idleTimeout: PropTypes.number.isRequired,
};

const mapStateToProps = state => {
  return {
    isLoggedIn: state.user.isLoggedIn,
  };
};
export default connect(mapStateToProps)(LogoutOnInactivity);
