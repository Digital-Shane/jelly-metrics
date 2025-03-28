# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-03-28
### Updated
- Active stream and device calculations. 
  * `IsActive` appears to be true for any returned session, so we can't use that for metrics.
  * `PlayState` is always returned, so utilize `NowPlayingItem` instead.
- Token auth to use supported query parameter instead of deprecated header. 

## [1.0.0] - 2025-03-27
### Released
- Metrics api to export media server metrics from the jellyfin API.
