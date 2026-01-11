package app

import (
	"os"

	"hytale-launcher/internal/buildscan"
	"hytale-launcher/internal/fork"
	"hytale-launcher/internal/hytale"
	"hytale-launcher/internal/net"
)

// PopulateSentryExtra populates the Sentry extra context with application state.
// This is called before sending error reports to provide debugging context.
func (a *App) PopulateSentryExtra(extra map[string]any) {
	// Add account information if logged in.
	if a.Auth != nil {
		if acct := a.Auth.GetAccount(); acct != nil {
			accountInfo := map[string]any{
				"profiles":         acct.Profiles,
				"patchlines":       acct.Patchlines,
				"selected_profile": acct.SelectedProfile,
				"selected_channel": acct.SelectedChannel,
			}
			extra["account"] = accountInfo
		}
	}

	// Add environment/channel state if available.
	if a.State != nil {
		envInfo := map[string]any{
			"dependencies":  a.State.Dependencies,
			"offline_ready": a.State.OfflineReady,
			"channel":       a.State.Channel,
		}
		extra["environment"] = envInfo
	}

	// Add launcher-specific information.
	launcherInfo := map[string]any{
		"net_mode":    net.Current(),
		"storage_dir": hytale.StorageDir(),
	}
	extra["launcher"] = launcherInfo

	// Add information about installed games.
	installs := buildscan.ScanInstalledGames(false)
	extra["installs"] = installs

	// Add process information.
	processInfo := map[string]any{
		"elevated": fork.IsElevated(),
		"pid":      os.Getpid(),
	}
	extra["process"] = processInfo
}
