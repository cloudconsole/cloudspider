// main.js
var React = require('react');
var ReactDOM = require('react-dom');

class Avatar extends React.Component {
    render() {
        return (
            <div className="row">
                <div className="col-md-10 col-md-offset-1">
                    <PagePic uriname={this.props.pagename}/>
                    <PageLink uriname={this.props.pagename}/>
                </div>
            </div>
        );
    }
}

class PagePic extends React.Component {
    render() {
        return (
            <img src={'https://graph.facebook.com/' + this.props.uriname + '/picture'}/>
        );
    }
}

class PageLink extends React.Component {
    render() {
        return (
            <a href={'https://www.facebook.com/' + this.props.uriname}>
                {this.props.uriname}
            </a>
        );
    }
}

ReactDOM.render(
    <Avatar pagename="Engineering"/>,
    document.getElementById('dashboard')
);

var myDivElement = <div className="row">
    <div className="col-md-10 col-md-offset-1">
        <h1 className="page-header">Work In Progress</h1>
    </div> </div>;

ReactDOM.render(
    myDivElement,
    document.getElementById('dashboard1')
);
