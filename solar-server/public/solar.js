/** @jsx React.DOM */

var toRGB = function(c) {
  return "rgb(" + c.R + "," + c.G + "," + c.B + ")"
}

var Timeline = React.createClass({displayName: "Timeline",

  getInitialState: function() {
    return {
      index: 0,
      avgColor: {
        R: 255,
        G: 255,
        B: 255
      },
      analyzerData: []
    };
  },

  componentDidMount: function() {
    this.getAnalyzerData();
    $(document).keydown($.proxy(function(e) {
    switch (e.which) {
          case 37: // left
              this.handleIndexDecrease();
              break;

          case 39: // right
              this.handleIndexIncrease();
              break;
      }
    }, this));
    window.onpopstate = $.proxy(this.handlePopState, this);
    var hash = window.location.hash;
    if (!!hash) {
      var i = parseInt(hash.substring(1));
      this.handleSetIndex(i);
    }
  },

  getAnalyzerData: function() {
    $.getJSON(this.props.root + "analysis", $.proxy(function(data) {
      this.setState({analyzerData: data, avgColor: data[this.state.index].AverageColor});
    }, this));
  },

  handlePopState: function(s) {
    this.handleSetIndex(s.state.index, true);
  },

  handleSetIndex: function(i, skipPushState) {
    this.setState({index: i});
    if (this.state.analyzerData.length > 0) {
      this.setState({avgColor: this.state.analyzerData[i].AverageColor});
    }

    if (!skipPushState) {
      window.history.pushState({index: i}, "Image " + i, "#" + i);
    }
  },

  handleIndexDecrease: function() {
    var index = this.state.index;
    if (index === 0) {
      index = this.state.analyzerData.length - 1;
    } else {
      index -= 1;
    }
    this.handleSetIndex(index);
  },

  handleIndexIncrease: function() {
    var index = this.state.index;
    index = (index + 1) % this.state.analyzerData.length;
    this.handleSetIndex(index);
  },

  render: function() {
    $("body").css("background-color", toRGB(this.state.avgColor));
    var index = this.state.index;
    return (
      React.createElement("div", {className: "timeline-container"}, 
        React.createElement(Image, {root: this.props.root, index: index, onIndexDecrease: this.handleIndexDecrease, onIndexIncrease: this.handleIndexIncrease, onSetIndex: this.handleSetIndex, analyzerData: this.state.analyzerData})
      )
    )
  }
});

var Image = React.createClass({displayName: "Image",

  getInitialState: function() {
    return {
      width: 1000,
      height: 562
    };
  },

  componentDidMount: function() {
    $(window).resize($.proxy(this.handleWinResize, this));
    this.handleWinResize();
  },

  handleWinResize: function() {
    var ow = 2592.0;
    var oh = 1458.0;
    var ww = $(window).width() - 50;
    var wh = $(window).height() - 50;

    var ra = ww / ow;
    var rb = wh / oh;

    var a = {
        w: ow * ra,
        h: oh * ra
    }

    var b = {
        w: ow * rb,
        h: oh * rb
    }

    var winner = a;

    if (a.w > ww || a.h > wh) {
        winner = b;
    }
    this.setState({width: winner.w, height: winner.h});
  },

  getImgSrc: function(idx) {
    if (this.props.analyzerData.length > 0) {
      return this.props.root + "images/perm/" + this.props.analyzerData[idx].Name;
    }
    return "";
  },

  handleClick: function(e) {
    var ref = "image" + this.props.index;
    var rect = $(this.refs[ref].getDOMNode())[0].getBoundingClientRect();
    var mid = (rect.right - rect.left) / 2.0;
    var offsetX = e.clientX - rect.left;
    if (offsetX < mid) {
      this.handleLeft();
    } else {
      this.handleRight();
    }
  },

  handleLeft: function() {
    this.props.onIndexDecrease();
  },

  handleRight: function() {
    this.props.onIndexIncrease();
  },

  renderImg: function(idx, first) {
    var style = {
      zIndex: idx === this.props.index ? 2 : 1
    };

    var className = "solar-image";
    if (first) {
      className = " solar-image-first"
    }

    return (
      React.createElement("img", {style: style, ref: "image" + idx, onClick: this.handleClick, width: this.state.width, height: this.state.heigh, className: className, src: this.getImgSrc(idx)})
    )
  },

  render: function () {

    var imgNodes = [];
    if (this.state.width > 1000) {
      var i = -1;
      var loadSize = 10;
      var start = Math.max(this.props.index - loadSize, 0)
      var end = Math.min(this.props.index + loadSize, this.props.analyzerData.length)
      for (var i = start; i < end; i++) {
        imgNodes.push(this.renderImg(i, i === start));
      }
    } else {
      imgNodes = [this.renderImg(this.props.index, true)]
    }

    return (
      React.createElement("div", null, 
        imgNodes, 
        React.createElement(ColorLine, {data: this.props.analyzerData, onSetIndex: this.props.onSetIndex}), 
        React.createElement("div", {className: "control-container"}, 
          React.createElement("span", {className: "left", onClick: this.handleLeft}, " ←"), 
          React.createElement("div", {className: "center"}, this.props.index), 
          React.createElement("span", {className: "right", onClick: this.handleRight}, " →")
        )
      )
    );
  }
});

var ColorLine = React.createClass({displayName: "ColorLine",

  render: function() {
    var widthp = (1/this.props.data.length) * 100;
    var i = -1;
    var nodes = this.props.data.map(function(d) {
      i++;
      return (
        React.createElement(ColorSegment, {index: i, widthp: widthp, data: d, onSetIndex: this.props.onSetIndex})
      );
    }, this);

    return (
      React.createElement("div", {id: "colorline"}, 
        nodes
      )
    );
  }
});

var ColorSegment = React.createClass({displayName: "ColorSegment",
  handleClick: function(e) {
    e.preventDefault();
    this.props.onSetIndex(this.props.index);
  },
  render: function() {
    var divStyle = {
      "background-color": toRGB(this.props.data.AverageColor),
      "width": this.props.widthp + "%"
    };
    return (
      React.createElement("div", {className: "color-segment", style: divStyle, onClick: this.handleClick})
    )
  }
});

var getRootPath = function() {
  var url = document.URL;
  var start = url.indexOf("//");
  var rootEnd = url.lastIndexOf('/');
  var root = url.substring(start, rootEnd + 1);
  return root;
};

React.render(
  React.createElement(Timeline, {root: getRootPath()}),
  document.getElementsByClassName("main-wrapper")[0]
);
